from typing import List

from bs4 import BeautifulSoup


class HtmlParser:
    """封装 BeautifulSoup 的简单 HTML 解析器。"""

    def get_links(self, html: str) -> List[str]:
        soup = BeautifulSoup(html, "html.parser")
        links: List[str] = []
        for a in soup.find_all("a", href=True):
            links.append(a["href"])
        return links

    def get_title(self, html: str) -> str:
        soup = BeautifulSoup(html, "html.parser")
        if soup.title and soup.title.string:
            return soup.title.string.strip()
        return ""

