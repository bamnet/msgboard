var msgboardControllers = angular.module('msgboardControllers', []);

msgboardControllers.controller('DisplayShowCtrl', ['$scope', '$interval', 'Blurbs', 'Page',
	function ($scope, $interval, Blurbs, Page) {
		$scope.activePageIndex = 0;
		$scope.activePage = undefined;
		$scope.pages = Page.list({view: 'ids'}, function(){
			$scope.activePage = Page.get({pageId: $scope.pages[$scope.activePageIndex].id});
		});
		$scope.blurbs = Blurbs.get();

		var pageInterval = $interval(function(){
			if($scope.activePageIndex >= $scope.pages.length-1) {
				Page.list({view: 'ids'}, function(pages){
					$scope.pages = pages;
					$scope.activePageIndex = 0;
					Page.get({pageId: $scope.pages[$scope.activePageIndex].id}, function(page) {
						$scope.activePage = page;
					});
				});
			} else {
				$scope.activePageIndex++;
				Page.get({pageId: $scope.pages[$scope.activePageIndex].id}, function(page) {
					$scope.activePage = page;
				});
			}
		}, 1000);

		var blurbsInterval = $interval(function(){
			Blurbs.get({}, function(blurbs){
				$scope.blurbs = blurbs;
			});
		}, 1000);

		$scope.stopIntervals = function(){
			var intervals = [pageInterval, blurbsInterval];
			angular.forEach(intervals, function(interval) {
				if (angular.isDefined(interval)) {
					$interval.cancel(interval);
					interval = undefined;
				}
			});
		};

		$scope.$on('$destroy', function() {
			// Stop any active intervals.
			$scope.stopIntervals();
		});
	}
]);
