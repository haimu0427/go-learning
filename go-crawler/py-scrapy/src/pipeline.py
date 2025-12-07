from typing import Any


class PrintPipeline:
    """示例数据管道：把抓取结果打印出来。"""

    def process_item(self, item: Any) -> None:
        print("[ITEM]", item)

