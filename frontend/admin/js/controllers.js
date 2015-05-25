var msgboardControllers = angular.module('msgboardControllers', []);

msgboardControllers.controller('PageListCtrl', ['$scope', 'Page',
	function ($scope, Page) {
		$scope.pages = Page.list();
	}
]);

msgboardControllers.controller('PageShowCtrl', ['$scope', '$routeParams', 'Page',
	function($scope, $routeParams, Page) {
		$scope.page =  Page.get({pageId: $routeParams.pageId});
	}
]);

msgboardControllers.controller('PageEditCtrl', ['$scope', '$routeParams', '$location', 'Page',
	function($scope, $routeParams, $location, Page) {
		var pageId = $routeParams.pageId;
		$scope.page = Page.get({pageId: pageId});

		$scope.update = function(page) {
			$scope.page = angular.copy(page);
			$scope.page.$update({pageId: pageId}, function(){
				$location.path('/pages/' + pageId);
			}, function(err){
				$scope.error = err.data;
			});
		};

		$scope.delete = function(page) {
			$scope.page.$delete({pageId: pageId}, function(){
				$location.path('/pages/');
			}, function(err){
				$scope.error = err.data;
			});
		};
	}
]);

msgboardControllers.controller('PageCreateCtrl', ['$scope', '$location', 'Page',
	function($scope, $location, Page) {
		$scope.page =  new Page();
		$scope.create = function(page) {
			$scope.page = angular.copy(page);
			$scope.page.$create(function(page){
				$location.path('/pages/' + page.id);
			}, function(err){
				$scope.error = err.data;
			});
		};
	}
]);

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
