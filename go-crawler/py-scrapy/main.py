"""简单命令行入口：运行示例爬虫。"""

from src.spider import ExampleSpider


def main() -> None:
    spider = ExampleSpider()
    spider.run()


if __name__ == "__main__":
    main()

