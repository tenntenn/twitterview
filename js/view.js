define(
    'view',
    [],
    function() {
        var viewId = $('body').data('view-id');
        var render = function(pdf, pageNum, scale) {
            pdf.getPage(pageNum).then(function(page) {
                var viewport = page.getViewport(scale);
                var canvas = $('#pdf-canvas')[0]
                var context = canvas.getContext('2d');
                canvas.height = viewport.height;
                canvas.width = viewport.width; 
                var renderContext = {
                    canvasContext: context,
                    viewport: viewport
                }
                page.render(renderContext);
            });
        };

        var pdfFile = "pdf/"+viewId+".pdf";
        PDFJS.getDocument(pdfFile).then(function(pdf) {
            var pageNum = 1;
            // render first page
            render(pageNum); 

            $('html').keydown(function(e) {
            });
        });
    }
);
