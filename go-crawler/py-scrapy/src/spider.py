from abc import ABC, abstractmethod
from typing import Iterable, Any, Optional

from .fetcher import HttpFetcher
from .parser import HtmlParser
from .pipeline import PrintPipeline


class BaseSpider(ABC):
    """爬虫基类：定义通用结构和流程。"""

    name: str = "base_spider"

    def __init__(self) -> None:
        self.fetcher = HttpFetcher()
        self.parser = HtmlParser()
        self.pipeline = PrintPipeline()

    @abstractmethod
    def start_urls(self) -> Iterable[str]:
        """起始 URL 列表。"""
        raise NotImplementedError

    @abstractmethod
    def parse(self, url: str, html: str) -> Iterable[Any]:
        """解析页面，产生 item（可以是 dict 或自定义对象）。"""
        raise NotImplementedError

    def run(self) -> None:
        for url in self.start_urls():
            resp = self.fetcher.get(url)
            if not resp:
                print(f"[ERROR] 请求失败: {url}")
                continue
            html = resp.text
            for item in self.parse(url, html):
                self.pipeline.process_item(item)


class ExampleSpider(BaseSpider):
    """一个简单示例：抓取页面标题。"""

    name = "example_spider"

    def __init__(self, urls: Optional[Iterable[str]] = None) -> None:
        super().__init__()
        self._urls = list(urls) if urls is not None else [
            "https://www.python.org/",
        ]

    def start_urls(self) -> Iterable[str]:
        return self._urls

    def parse(self, url: str, html: str):
        title = self.parser.get_title(html)
        yield {
            "url": url,
            "title": title,
        }

