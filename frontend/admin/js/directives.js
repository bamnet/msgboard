var msgboardDirectives = angular.module('msgboardDirectives', []);

msgboardDirectives.directive('richTextEditor', function() {
	return {
		restrict: 'A',
		require: 'ngModel',
		transclude: true,
		templateUrl: 'partials/directives/richtexteditor.html',
		link: function(scope, element, attrs, ctrl) {
			var toolbarElem = element.children()[0];
			var editorElem = element.children()[1];
			var editor = new wysihtml5.Editor(editorElem, {
				toolbar: toolbarElem,
				parserRules: wysihtml5ParserRules
			});

			// Update the model on changes to the view.
			editor.on('change', function() {
				scope.$apply(function() {
					ctrl.$setViewValue(editor.getValue());
				});
			});

			// Load the model into the editor.
			ctrl.$render = function() {
				editor.setValue(ctrl.$viewValue);
			};
			ctrl.$render();
		}
	};
});
