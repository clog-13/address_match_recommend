# 使用说明

1. 导入`./resource/testdb_v0.sql`，国家标准4级行政区域库（可提升到5级） 和 预存地址数据
2. （可选）导入默认预存地址数据，`test_addresses.txt`， 执行`./sh/import_base_addrs.go`导入
3. 查询接口`core.FindsimilarAddress(addressText string, topN int, explain bool) (Query, bool)` 
4. 添加数据接口`core.ImportAddr(text string)`


# 算法说明
参照
https://github.com/liuzhibin-cn/address-semantic-search
1. 将 地址字符串 解析为Address

   ```go
   type Address struct {
       Id int64 `gorm:"primaryKey;comment:地址ID" json:"ID"`
   
       RawText     string `gorm:"type:text;" json:"raw_text"`
       AddressText string `gorm:"type:text;" json:"address_text"`
       RoadText    string `gorm:"type:text;" json:"road"`
       RoadNum     string `gorm:"type:text;" json:"road_num"`
       BuildingNum string `gorm:"type:text;" json:"building_num"`
   
       ProvinceId, CityId, DistrictId, StreetId, VillageId, TownId uint
   
       Province *Region `gorm:"-"`
       City     *Region `gorm:"-"`
       District *Region `gorm:"-"`
       Street   *Region `gorm:"-"`
       Town     *Region `gorm:"-"`
       Village  *Region `gorm:"-"`
   }
   type Region struct {
       ID uint `gorm:"primaryKey;comment:行政区域ID" json:"ID"`
   
       ParentID uint   `gorm:"type:uint;" json:"region_parent_id"`
       Name     string `gorm:"type:string;" json:"region_name"`
       Alias    string `gorm:"type:string;" json:"region_alias"`
       Types    int    `gorm:"type:SMALLINT;" json:"region_types"`
   
       Children     []*Region      `gorm:"-"`
       OrderedNames pq.StringArray `gorm:"-"`
       //_varchar OrderedNames pq.StringArray `gorm:"type:varchar(255)[]" json:"region_ordered_names"`
   }
   ```

2. 将 Address 解析为 Document

   ```go
   type Document struct {
       Id uint
   
       // 文档所有词条, 按照文档顺序, 未去重
       Terms    []*Term
       TermsMap map[string]*Term
   
       TownId       uint
       Town         *Term // 乡镇相关的词条信息
       VillageId    uint
       Village      *Term
       RoadId       uint
       Road         *Term // 道路信息
       RoadNumId    uint
       RoadNum      *Term
       RoadNumValue int
   }
   type Term struct {
       Id uint
   
       TermId uint
       Text   string
       Types  int
       Idf    float64
   
       Ref *Term
   }
   ```

3. 利用 文本TF-IDF余弦相似度算法 和 Lucene的评分算法 计算相似度

   **TF-IDF余弦相似度：**

   - `TC`: `Term Count`，词数，某个词在文档中出现的次数。
   - `TF`: `Term Frequency`，词频，某个词在文档中出现的频率，`TF = 该词在文档中出现的次数 / 该文档的总词数`。
   - `IDF`: `Inverse Document Frequency`，逆文档词频，`IDF = log( 文档总数 / ( 包含该词的文档数 + 1 ) )`。分母加1是为了防止分母出现0的情况。
   - `TF-IDF`: 词条的特征值，`TF-IDF = TF * IDF`。

   **Lucene的评分算法**
   
   [![Lucene评分算法](https://camo.githubusercontent.com/c535782e472b53437d49eece055e69caea25e007bd479e0ebc7d6ac9dd3b9fae/68747470733a2f2f7269636869652d6c656f2e6769746875622e696f2f79647265732f696d672f31302f3138302f313031342f6c7563656e652d73636f72652d66756e6374696f6e2e706e67)](https://camo.githubusercontent.com/c535782e472b53437d49eece055e69caea25e007bd479e0ebc7d6ac9dd3b9fae/68747470733a2f2f7269636869652d6c656f2e6769746875622e696f2f79647265732f696d672f31302f3138302f313031342f6c7563656e652d73636f72652d66756e6374696f6e2e706e67)
   
   - **① score(q, d)**
     查询文档q与文档d的相关性评分值，Lucene中的评分值是一个大于0的实数值。
   
   - **② queuryNorm(q)**
     `queryNorm(q) = 1 / sqrt( sumOfSquaredWeights )`，`sumOfSquaredWeights`是查询文档q中每个词条的`IDF`平方根之和。同理，文档库中的每个文档也会使用这一公式进行归一化处理。
     `queryNorm`的目的是将lucene评分值进行归一化处理，使不同文档之间的评分值具有可比较性，但这个归一化算法严谨性有待证明。相比较之下，余弦相似度算法的结果无需额外的归一化处理。
   
   - **③ coord(q, d)**
     `coord(q, d) = 文档d中匹配上的词条数量 / 文档q的词条数量`
     `coord`的效果是根据匹配词条数量进行加权，匹配词条越多加权值越高，表示查询文档q与文档d相关性越高，如果查询文档q的所有词条都能匹配，则`coord`值为1。
   
   - **④ ∑(...)(t in q)**
     对查询文档q中的每个词条t，使用公式⑤⑥⑦⑧求值，将得到的值求和汇总。
   
   - **⑤ tf(t in d)**
     注意，这里的TF是词条t在文档d中的词频。
   
   - **⑥ idf(t)²**
   
   - ⑦ t.getBoost()
   
     Lucene的数据结构：
   
     - 每个文档Document由N个Field组成，Field类型可以是数值、文本、日期等，文本类型的Field经过分词后会包含N个词条Term。
     - 可以使用索引(Index)为文档分组，例如按时间方式建立文档索引，索引`2016_09`中存放9月份的文档，索引`2016_10`中存放10月份的文档。 Lucene提供`Boost`参数为相关性加权，可以为文档的不同Field设置不同的`Boost`，也可以为索引设置`Boost`（例如索引`2016_10` `Boost`=2.0，索引`2016_09` `Boost`=1.5）。
   
   - **⑧ norm(t,d)**
     `norm(t, d) = 1 / sqrt( numTerms )`，`numTerms`是文档d的词条数量。
     `norm(t, d)`是基于文档/Field长度进行的归一化处理，其效果是，同一个词条，出现在较短的Field（例如Title）中，比出现在较长的Field（例如Content）的相关性更高。
     另外，`TF-IDF`的一些实际运用中会对TF使用这种方式进行计算，即`TF = 该词在文档中出现的次数 / sqrt( 该文档的总词数 )`

4. 查询前先进行布隆过滤器判断

# TODO

- 代码健壮性，数据库操作（存在则不更新...） 



- 优化 地址分解细节 和 正则表达式

- 增加语料库（“公馆”，“乐园”，“雅居”，...）， 提升解析细度

  

- 线程安全，datarace

- 减少空间复杂度（优化内存数据结构中的冗余数据，...）

- 添加缓存

- ...

# 已知部分问题
- 解析不完全 导致相同字符串相似性计算不能到 1
- gse分词库 会丢弃 ‘-’ 后的内容（“xxxx小区2-1-17” -> “xxxx小区, 2”），暂时处理为“xxxx小区2117”
- 内存空间数据浪费