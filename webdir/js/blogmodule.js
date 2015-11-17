var app = angular.module("blog", ['ui.bootstrap','btford.markdown'])
.config(['markdownConverterProvider', function (markdownConverterProvider) {
  // options to be passed to Showdown
  // see: https://github.com/coreyti/showdown#extensions
  //markdownConverterProvider.config({
  //  extensions: ['youtube']
  //});
}])
