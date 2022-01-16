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
    <style>
        .modal {
            display: none; /* Hidden by default */
            position: fixed; /* Stay in place */
            z-index: 1000; /* Sit on top */
            padding-top: 100px; /* Location of the box */
            left: 0;
            top: 0;
            width: 100%; /* Full width */
            height: 100%; /* Full height */
            overflow: auto; /* Enable scroll if needed */
            background-color: rgb(0, 0, 0); /* Fallback color */
            background-color: rgba(0, 0, 0, 0.9); /* Black w/ opacity */
        }

        /* Modal Content (Image) */
        .modal-content {
            margin: auto;
            display: block;
        }

        /* Add Animation - Zoom in the Modal */
        .modal-content, #caption {
            animation-name: zoom;
            animation-duration: 0.6s;
        }

        @keyframes zoom {
            from {
                transform: scale(0)
            }
            to {
                transform: scale(1)
            }
        }
    </style>
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
<div id="myModal" class="modal" onclick="disShowImage()">
    <img class="modal-content" id="img01" onmousewheel="return zoomImg(this)"
         style="width: auto;height: auto;max-width: 100%;max-height: 100%;">
</div>
<script type="text/javascript">
`
	PostPageFooter = `</script>
<script type="text/javascript">
	var modal = document.getElementById("myModal");
	var modalImg = document.getElementById("img01");

	// When the user clicks on <span> (x), close the modal
	function disShowImage() {
		modal.style.display = "none";
	}

	function zoomImg(obj) {
		// 一开始默认是100%
		let zoom = parseInt(obj.style.zoom, 10) || 100;
		// 滚轮滚一下wheelDelta的值增加或减少120
		zoom += event.wheelDelta / 12;
		if (zoom > 0) {
			obj.style.zoom = zoom + '%';
		}
		return false;
	}
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
					$("#preview img").on('click',function(){
                        console.log("test---",$(this).attr("src"))
                        modal.style.display = "block";
                        modalImg.style.zoom = "reset";
                        modalImg.src = $(this).attr("src");
                        imageClicked = true;
                    })
				},
			})
		document.onkeydown = function (oEvent) {
            var oEvent = oEvent || window.oEvent;
            //获取键盘的keyCode值
            var nKeyCode = oEvent.code;
            //获取ctrl 键对应的事件属性
            var bCtrlKeyCode = oEvent.ctrlKey || oEvent.metaKey;
            if (nKeyCode === "KeyS" && bCtrlKeyCode) {
                console.log("ctrl + s");
            } else if (nKeyCode === "Escape" && imageClicked) {
                imageClicked = false;
                modal.style.display = "none";
            }
        }
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
