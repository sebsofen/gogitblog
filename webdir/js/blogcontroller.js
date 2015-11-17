app.controller("blogController", function($scope, $location, $http) {
    $scope.post_slug = "John";
    $scope.lastName= "Doe";

    $scope.location = $location;
    $scope.$watch('location.search()', function() {
        $scope.post_slug = ($location.search()).post_slug;

        $http.get('/posts/' + ($location.search()).post_slug).
                then(function(data) {
                  console.log(data.data.Post)
                    $scope.post = data.data.Post;
                });

    }, true);


});
