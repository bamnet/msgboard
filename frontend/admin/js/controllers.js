var msgboardControllers = angular.module('msgboardControllers', []);

msgboardControllers.controller('HeaderController', ['$scope', '$location',
	function($scope, $location){
		$scope.isActive = function(name) {
			return $location.path().indexOf('/' + name) >= 0;
		};
	}
]);

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

msgboardControllers.controller('BlurbsShowCtrl', ['$scope', 'Blurbs',
	function($scope, Blurbs) {
		$scope.blurbs =  Blurbs.get();
	}
]);

msgboardControllers.controller('BlurbsEditCtrl', ['$scope', '$location', 'Blurbs',
	function($scope, $location, Blurbs) {
		$scope.blurbs = Blurbs.get();

		$scope.update = function(blurbs) {
			$scope.blurbs = angular.copy(blurbs);
			$scope.blurbs.$update(null, function(){
				$location.path('/blurbs');
			}, function(err){
				$scope.error = err.data;
			});
		};
	}
]);
