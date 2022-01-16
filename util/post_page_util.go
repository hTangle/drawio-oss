package util

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
)

const (
	PostPageHeader = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8"/>
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
    <title>SuperMarkdownEditor</title>

    <script src="https://cdn.bootcdn.net/ajax/libs/jquery/3.6.0/jquery.min.js"></script>
    <link href="https://cdn.bootcdn.net/ajax/libs/twitter-bootstrap/4.6.0/css/bootstrap.min.css" rel="stylesheet">
    <script src="https://cdn.bootcdn.net/ajax/libs/twitter-bootstrap/4.6.0/js/bootstrap.min.js"></script>
    <script src="https://cdn.bootcdn.net/ajax/libs/twitter-bootstrap/4.6.0/js/bootstrap.bundle.min.js"></script>
    <link href="https://cdn.bootcdn.net/ajax/libs/font-awesome/5.15.1/css/fontawesome.min.css" rel="stylesheet">

    <link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/vditor/dist/index.css"/>
    <script src="https://cdn.jsdelivr.net/npm/vditor/dist/method.min.js"></script>
    <link href="/css/sticky-footer-navbar.css" rel="stylesheet">
</head>

<body>
<header>
    <!-- Fixed navbar -->
    <nav class="navbar navbar-expand-md navbar-dark fixed-top bg-dark" style="padding-bottom: 0.2rem">
        <div class="container-fluid">
            <a class="navbar-brand" href="#">ahsup</a>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarCollapse"
                    aria-controls="navbarCollapse" aria-expanded="false" aria-label="Toggle navigation">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navbarCollapse">
                <ul class="navbar-nav me-auto mb-2 mb-md-0">
                    <li class="nav-item">
                        <a class="nav-link active" aria-current="page" href="/gallery">Home</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="/main">Edit</a>
                    </li>
                    <li class="nav-item">
                        <a class="nav-link" href="/login">Login</a>
                    </li>
                </ul>
            </div>
        </div>
    </nav>
</header>
<div class="container-fluid h-100" style="width: 75%">
    <h1 id="articleTitle" style="text-align: center;margin-top: 2em">Markdown转HTML的显示处理之自定义 ToC 容器</h1>
    <div id="outline"></div>
    <div id="preview" class="vditor" style="padding-right: 250px;"></div>
</div>
<footer class="footer mt-auto py-3 bg-light">
    <div class="container">
        <span class="text-muted" style="text-align: center;"><p class="text-muted">&copy; 2021–2022 皖ICP备2022000174号</p></span>
    </div>
</footer>
<script type="text/javascript">
`
	PostPageFooter = `</script>
<script type="text/javascript">
    const initOutline = () => {
        const headingElements = []
        Array.from(document.getElementById('preview').children).forEach((item) => {
            if (item.tagName.length === 2 && item.tagName !== 'HR' && item.tagName.indexOf('H') === 0) {
                headingElements.push(item)
            }
        })

        let toc = []
        window.addEventListener('scroll', () => {
            const scrollTop = window.scrollY
            toc = []
            headingElements.forEach((item) => {
                toc.push({
                    id: item.id,
                    offsetTop: item.offsetTop,
                })
            })

            const currentElement = document.querySelector('.vditor-outline__item--current')
            for (let i = 0, iMax = toc.length; i < iMax; i++) {
                if (scrollTop < toc[i].offsetTop - 30) {
                    if (currentElement) {
                        currentElement.classList.remove('vditor-outline__item--current')
                    }
                    let index = i > 0 ? i - 1 : 0
                    document.querySelector('span[data-target-id="' + toc[index].id + '"]').classList.add('vditor-outline__item--current')
                    break
                }
            }
        })
    }

    $(function () {
		$("#articleTitle").text(title);
		$(document).attr("title",title);
		Vditor.preview(document.getElementById('preview'),
			content, {
				anchor: 1,
				after () {
					if (window.innerWidth <= 768) {
						return
					}
					const outlineElement = document.getElementById('outline')
					Vditor.outlineRender(document.getElementById('preview'), outlineElement)
					if (outlineElement.innerText.trim() !== '') {
						outlineElement.style.display = 'block';
						// outlineElement.style.float = 'right';
						outlineElement.style.position='fixed';
						outlineElement.style.top='10em';
						outlineElement.style.right='0px';
						initOutline()
					}
				},
			})
    });
</script>
</body>
<body>
</body>
</html>`
	TitleHtmlTemplate   = "var title=`%s`;\n"
	ContentHtmlTemplate = "var content=`%s`;"
)

func SavePostToHTML(id, title, html, targetPath string) {

	AllHtml := PostPageHeader + fmt.Sprintf(TitleHtmlTemplate, title) + fmt.Sprintf(ContentHtmlTemplate, strings.ReplaceAll(html, "`", "\\`")) + PostPageFooter
	ioutil.WriteFile(path.Join(targetPath, id+".html"), []byte(AllHtml), 0666)
}
