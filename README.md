<!--
 * @Author: qizk qizk@mail.open.com.cn
 * @Date: 2024-06-28 14:23:49
 * @LastEditors: qizk qizk@mail.open.com.cn
 * @LastEditTime: 2024-08-21 17:56:47
 * @FilePath: \word2excel\README.md
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
-->

# word2excel

试题模板转换工具，该项目根据业务需要会将指定格式的模板转换为Excel格式的导题模板。

目前支持的转化模板分为**两种**, **Word**和**Html**格式。

1. 请将Word格式的模板文件放置于**runtime/word**目录下，目前仅支持三种格式的word文件，三种格式的word文件分别对应根目录下的**word_template1.docx**, **word_template2.docx**, **word_template3.docx**， 符合这三种格式的word文件可以被正确解析， 解析后的Excel格式的试题导入模板会被存放在**runtime/excel**目录下。
2. 请将Html格式的模板文件放置于**runtime/html**目录下，目前仅支持一种格式的html文件，格式文件对应根目录下的**html_template.html**， 解析后的Excel格式的试题导入模板会被存放在**runtime/excel**目录下。

## 执行命令如下

```bash
go run . -n 2 -t one
```

> 注： -n 控制并发数量 -t 执行的模板类型，目前仅支持四种类型，分别对应的参数为：one, two, three，html。html对应解析html格式的文件，其他对应相应的word文件
