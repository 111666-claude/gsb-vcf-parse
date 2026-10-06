# vcf-parse

解析 vCard 文本：折行还原、参数与值拆分、值里的转义还原，并按姓名去重建索引。只用标准库。

```
go test ./...
go vet ./...
go run ./cmd/vcfparse --sample=fold
go run ./cmd/vcfparse --sample=esc
go run ./cmd/vcfparse --sample=param
go run ./cmd/vcfparse --sample=dup
go run ./cmd/vcfparse --sample=work
```

## 口径

- **折行**：以空格或制表符开头的行是上一条的续行，去掉行首空白后直接接上。
- **参数**：`名字;参数=值:值` 里，第一个分号之前是属性名，分号之后到冒号之前是参数。
- **转义**：值里的 `\,` `\;` `\n` `\\` 还原成对应单字符（换行、反斜杠等）。
- **去重**：同一个 `FN` 的卡片只保留第一次出现的那张。
- **查找**：按 `FN` 查找走索引，不逐张比较。

## 不变量

- 每张卡片的属性名不含分号，值里不含转义用的反斜杠。
- 同一个 `FN` 在通讯录里只有一张卡片。
- `scanned` 不随卡片数乘查询次数放大：八百张卡片查八百次的 `scanned` 不超过 6000。

## 输出契约

`Card` 是属性列表；`Get(name)` 取第一个同名属性的值；`Find(fn)` 按姓名取卡片。
场景打印一行 JSON。`--sample=work` 打印 `{"cards": N, "scanned": N}`。
