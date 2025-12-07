import time
from typing import Optional

import requests


class HttpFetcher:
    """简单的 HTTP 请求封装，带重试和超时。"""

    def __init__(self, timeout: int = 10, max_retries: int = 3, retry_delay: float = 1.0):
        self.timeout = timeout
        self.max_retries = max_retries
        self.retry_delay = retry_delay

    def get(self, url: str, **kwargs) -> Optional[requests.Response]:
        headers = kwargs.pop("headers", {
            "User-Agent": "SimpleCrawler/0.1 (+learning)"
        })
        for attempt in range(1, self.max_retries + 1):
            try:
                resp = requests.get(url, headers=headers, timeout=self.timeout, **kwargs)
                resp.raise_for_status()
                return resp
            except requests.RequestException:
                if attempt == self.max_retries:
                    return None
                time.sleep(self.retry_delay)

