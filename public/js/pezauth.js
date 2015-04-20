var pezAuth = angular.module('pezAuth', [], function($interpolateProvider) {
      $interpolateProvider.startSymbol('{*{');
      $interpolateProvider.endSymbol('}*}');
  })
  .controller('PezAuthController', function($scope, $http, $timeout) {
    var myData = {};
    var pauth = this;
    var restUriBase = "/v1/auth/api-key";
    var username = "jcalabrese@pivotal.io";
    var restUri = [restUriBase, username].join("/")

    $timeout(function () {  
      callUsingVerb($http.get);
    }, 100);

   
    pauth.create = function() {
      if ( myData.ApiKey === "" || angular.isUndefined(myData.ApiKey) ) {
        callUsingVerb($http.put)

      } else {
        callUsingVerb($http.post)
      }
    };
 
    pauth.remove = function() {
      callUsingVerb($http.delete)
    };   

    function callUsingVerb(verbCaller) {
      var responsePromise = verbCaller(restUri);
      responsePromise.success(function(data, status, headers, config) {
          console.log(angular.toJson(data, true));
          $scope.myData = data;
      });
      
      responsePromise.error(function(data, status, headers, config) {
          alert("AJAX failed!");
      });
    }
  });
