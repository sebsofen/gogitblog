app.controller("blogController", function($scope, $location, $http) {
    $scope.post_slug = "John";
    $scope.lastName= "Doe";
    $scope.posts = Array();
    $scope.location = $location;
    $scope.offset = 0;
    $scope.$watch('location.search()', function() {
        $scope.post_slug = ($location.search()).post_slug;

        $http.get('/posts/' + ($location.search()).post_slug).
                then(function(data) {
                  console.log(data.data.Post)
                    $scope.post = data.data.Post;
                });

    }, true);


    $scope.load_more = function() {
      $http.get('/listposts/' + $scope.offset + '/10').then(
        function(data) {

          for(i = 0; i < data.data.length; i++){
            console.log(data.data[i]);
            $scope.posts.push(data.data[i]);
            
          }
        }
      );
      $scope.offset += 10;
    }


});
