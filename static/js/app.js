$(function() {

	imagesLoaded('.grid', function(grid) {
		var msnry = new Masonry('.grid', {itemSelector: '.grid-item', columnWidth: 1});
	});	

	$('.description').on('show.bs.collapse', function () {
		// scroll to top to see content
	  	$('body').animate({scrollTop: 0},'slow');
	});

});

