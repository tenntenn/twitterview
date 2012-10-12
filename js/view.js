define(
    'view',
    [],
    function() {

        // id this view
        var viewId = $('body').data('view-id');

        // pdf file name
        var pdfFile = "pdf/"+viewId+".pdf";

        // render a pdf page with given scale
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

        // create key event
        var createKeyEvent = function(pdf) {
            return function(e) {
                switch(e.which) {
                // next page
                case 13: // enter key
                case 39: // right arrow
                    break;
                // previous page
                case 8:
                case 37:
                    pageNum--;
                    if (pageNum <= 0) {
                       pageNum = 0; 
                    }
                    render(pdf, pageNum, scale);
                    break;
                }
            };
        };

        // load and render pdf document 
        PDFJS.getDocument(pdfFile).then(function(pdf) {
            var pageNum = 1;
            var scale = 1;

            // render first page
            render(pageNum); 

            // regist key event
            $('html').keydown(createKeyEvent(pdf));
        });

        // Tweets
        var tweets = ko.observableArray();
        ko.applyBindings(tweets, $('#tweets'));
    }
);
