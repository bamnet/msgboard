var msgboardServices = angular.module('msgboardServices', ['ngResource']);

msgboardServices.factory('Page', ['$resource',
	function($resource){
		return $resource('/api/pages/:pageId', {}, {
			list: {method:'GET', isArray:true},
			update: {method:'PATCH'},
			create: {method: 'POST'},
			delete: {method: 'DELETE'}
		});
	}
]);

msgboardServices.factory('Blurbs', ['$resource',
	function($resource){
		return $resource('/api/blurbs', {}, {
			get: {method:'GET'},
			update: {method:'PATCH'}
		});
	}
]);
