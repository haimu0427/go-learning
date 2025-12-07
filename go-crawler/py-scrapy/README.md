# Python 爬虫学习项目（轻量框架）

这是一个面向初学者的 Python 爬虫学习项目，不依赖 Scrapy，只用 `requests` + `BeautifulSoup`，代码结构清晰、便于扩展。

## 环境准备

1. 安装依赖（建议使用虚拟环境）：

```bash
cd d:/AAAmycode/go-learning/go-crawler/py-scrapy
python -m venv .venv
.venv/Scripts/activate  # Windows PowerShell 下可用: .venv/Scripts/Activate.ps1
pip install -r requirements.txt
```

> 如果你使用的是 PowerShell，可以这样激活虚拟环境：
>
> ```powershell
> cd d:/AAAmycode/go-learning/go-crawler/py-scrapy
> .venv/Scripts/Activate.ps1
> ```

## 项目结构

- `main.py`：命令行入口，运行示例爬虫
- `src/fetcher.py`：HTTP 请求封装（带重试和超时）
- `src/parser.py`：HTML 解析工具（基于 BeautifulSoup）
- `src/pipeline.py`：数据管道示例（打印 item）
- `src/spider.py`：爬虫基类和示例 `ExampleSpider`
- `tests/`：单元测试示例

## 运行示例爬虫

在项目根目录（包含 `main.py` 的目录）执行：

```bash
python main.py
```

程序会请求 `https://www.python.org/`，并打印出页面标题等信息。

## 运行测试

在项目根目录运行：

```bash
pytest
```

如果未安装 pytest，可以先：

```bash
pip install pytest
```

## 如何自己写一个爬虫

1. 在 `src/spider.py` 中参考 `ExampleSpider`：
   - 继承 `BaseSpider`
   - 实现 `start_urls`（返回要爬的 URL 列表）
   - 实现 `parse(self, url, html)`，解析页面并 `yield` 出 item（通常是 `dict`）
2. 在 `main.py` 中导入并运行你自己的 Spider。

示例（伪代码）：

```python
from src.spider import BaseSpider

class MySpider(BaseSpider):
    def start_urls(self):
        return ["https://example.com"]

    def parse(self, url, html):
        title = self.parser.get_title(html)
        yield {"url": url, "title": title}
```

## 学习建议

1. **第一步：掌握基础 HTTP & HTML**
   - 了解什么是 URL、请求方法（GET/POST）、状态码。
   - 学会用浏览器开发者工具查看网页结构（Elements / Network 面板）。

2. **第二步：熟悉 `requests` 的基本用法**
   - `requests.get(url)`，查看 `response.status_code`、`response.text`。
   - 设置请求头（如 User-Agent）、超时、简单错误处理。

3. **第三步：熟悉 BeautifulSoup**
   - 用 `BeautifulSoup(html, "html.parser")` 解析页面。
   - 学会 `find` / `find_all`、选择标签、读取 `text` 和属性（如 `href`）。
   - 在项目中多改一改 `HtmlParser` 的方法练习。

4. **第四步：理解爬虫框架中的几个关键概念**
   - **Fetcher**：负责下载网页（这里是 `HttpFetcher`）。
   - **Parser**：负责解析网页、提取数据（这里是 `HtmlParser`）。
   - **Spider**：组织抓取流程（起始 URL、如何解析每个页面）。
   - **Pipeline**：处理结果数据（保存到文件 / 数据库，这里先打印）。

5. **第五步：做两个小练习**
   - 写一个爬虫：抓取某个博客首页所有文章标题和链接。
   - 写一个爬虫：抓取某个电商网站某个搜索页的商品标题和价格（注意遵守网站的 robots 协议和服务条款）。

6. **第六步：进阶思路**
   - 增加：简单的 URL 去重、日志输出、随机 User-Agent。
   - 尝试：把结果保存到 CSV/JSON 文件。
   - 再进阶：了解更成熟的框架（如 Scrapy），思考和本项目的区别。

如果你愿意，我可以根据你想抓的具体网站，手把手带你写第一个实战爬虫。
