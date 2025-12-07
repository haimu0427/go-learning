from src.parser import HtmlParser


def test_get_title_and_links():
    html = """
    <html>
      <head><title>Test Page</title></head>
      <body>
        <a href="https://example.com">Example</a>
        <a href="/relative">Rel</a>
      </body>
    </html>
    """
    parser = HtmlParser()
    title = parser.get_title(html)
    links = parser.get_links(html)

    assert title == "Test Page"
    assert "https://example.com" in links
    assert "/relative" in links

