var msgboardControllers = angular.module('msgboardControllers', []);

msgboardControllers.controller('DisplayShowCtrl', ['$scope', '$interval', 'Page',
	function ($scope, $interval, Page) {
		$scope.activePageIndex = 0;
		$scope.activePage = undefined;
		$scope.pages = Page.list({view: 'ids'}, function(){
			$scope.activePage = Page.get({pageId: $scope.pages[$scope.activePageIndex].id});
		});

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

		$scope.stopIntervals = function(){
			if (angular.isDefined(pageInterval)) {
				$interval.cancel(pageInterval);
				pageInterval = undefined;
			}
		};

		$scope.$on('$destroy', function() {
			// Stop any active intervals.
			$scope.stopIntervals();
		});
	}
]);
