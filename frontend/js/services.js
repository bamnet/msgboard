var msgboardServices = angular.module('msgboardServices', ['ngResource']);

msgboardServices.factory('Page', ['$resource',
  function($resource){
    return $resource('pages/:pageId', {}, {
      list: {method:'GET', isArray:true}
    });
  }
]);