(function () {
    'use strict';
    var api = 'https://docs.wjdun.cn/api/articles';
    var category = 'cdn新闻';
    var container = document.querySelector('.news-content');
    if (!container) return;
    var params = new URLSearchParams(window.location.search);
    var id = params.get('id');
    var page = Math.max(1, parseInt(params.get('page'), 10) || 1);
    var size = 10;

    function element(tag, text) {
        var node = document.createElement(tag);
        if (text !== undefined) node.textContent = text;
        return node;
    }
    function link(text, query) {
        var node = element('a', text);
        node.href = 'news.html' + query;
        return node;
    }
    function message(text) {
        container.replaceChildren(element('p', text));
    }
    async function request(path) {
        var controller = new AbortController();
        var timer = setTimeout(function () { controller.abort(); }, 15000);
        try {
            var response = await fetch(api + path, { cache: 'no-store', credentials: 'omit', signal: controller.signal });
            if (!response.ok) throw new Error('HTTP ' + response.status);
            var result = await response.json();
            if (result.code !== 0 || !result.data) throw new Error('Invalid article response');
            return result.data;
        } finally { clearTimeout(timer); }
    }
    function date(value) {
        return typeof value === 'string' ? value.slice(0, 10) : '';
    }
    async function list() {
        var data = await request('?category=' + encodeURIComponent(category) + '&page=' + page + '&pageSize=' + size);
        var pages = Math.max(1, Math.ceil(data.total / size));
        if (page > pages) {
            window.location.replace('news.html?page=' + pages);
            return;
        }
        var ul = element('ul');
        ul.className = 'news-list';
        data.list.forEach(function (article) {
            var li = element('li');
            var heading = element('h3');
            heading.appendChild(link(article.title, '?id=' + encodeURIComponent(article.id)));
            li.append(heading, element('span', date(article.createdAt)));
            ul.appendChild(li);
        });
        if (!data.list.length) ul.appendChild(element('li', '暂无已发布的 CDN 新闻'));
        var paging = element('div');
        paging.className = 'paging';
        var navigation = element('div');
        navigation.className = 'pagination';
        navigation.appendChild(element('span', '当前为' + page + '/' + pages + '页 共' + data.total + '条记录 '));
        [['首页', 1], ['上一页', page - 1], ['下一页', page + 1], ['尾页', pages]].forEach(function (item, index) {
            var disabled = index < 2 ? page === 1 : page === pages;
            var node = disabled ? element('span', item[0]) : link(item[0], '?page=' + item[1]);
            node.style.margin = '0 6px';
            navigation.appendChild(node);
        });
        paging.appendChild(navigation);
        container.replaceChildren(ul, paging);
    }
    async function detail() {
        if (!/^[1-9]\d*$/.test(id)) throw new Error('Invalid article ID');
        var article = await request('/' + encodeURIComponent(id));
        if (article.category !== category) throw new Error('Invalid article category');
        var heading = element('h2', article.title);
        var meta = element('p', date(article.createdAt) + (article.author ? ' ' + article.author : ''));
        var body = element('div');
        body.className = 'news-article-body';
        var parsed = new DOMParser().parseFromString(article.content || '', 'text/html');
        var allowed = new Set(['P', 'DIV', 'SPAN', 'BR', 'HR', 'H1', 'H2', 'H3', 'H4', 'H5', 'H6', 'STRONG', 'B', 'EM', 'I', 'U', 'S', 'UL', 'OL', 'LI', 'BLOCKQUOTE', 'PRE', 'CODE', 'TABLE', 'THEAD', 'TBODY', 'TR', 'TH', 'TD', 'A', 'IMG']);
        Array.from(parsed.body.querySelectorAll('*')).forEach(function (node) {
            if (!allowed.has(node.tagName)) { node.remove(); return; }
            Array.from(node.attributes).forEach(function (attribute) {
                var name = attribute.name;
                var safe = ['href', 'src', 'alt', 'title', 'colspan', 'rowspan'].includes(name);
                var dataImage = node.tagName === 'IMG' && name === 'src' && /^data:image\/(?:png|jpe?g|gif|webp|avif|bmp);base64,[a-z0-9+/]+=*$/i.test(attribute.value);
                if ((name === 'href' || name === 'src') && !dataImage && !/^https?:$/i.test(new URL(attribute.value, api).protocol)) safe = false;
                if (!safe) node.removeAttribute(name);
            });
            ['href', 'src'].forEach(function (name) {
                if (node.hasAttribute(name)) node.setAttribute(name, new URL(node.getAttribute(name), api).href);
            });
            if (node.tagName === 'IMG') node.style.maxWidth = '100%';
        });
        body.append.apply(body, Array.from(parsed.body.childNodes));
        container.replaceChildren(link('返回新闻列表', ''), heading, meta, body);
        document.title = article.title + '_稳坚盾';
    }
    message('正在加载 CDN 新闻…');
    (id !== null ? detail() : list()).catch(function (error) {
        console.error('CDN news loading failed:', error);
        message(id !== null ? '文章不存在、未发布或暂时无法加载。' : '新闻暂时无法加载，请稍后重试。');
        container.appendChild(link(id !== null ? '返回新闻列表' : '重新加载', ''));
    });
})();
