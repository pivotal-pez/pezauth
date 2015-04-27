var pezAuth = angular.module('pezAuth', [], function($interpolateProvider) {
      $interpolateProvider.startSymbol('{*{');
      $interpolateProvider.endSymbol('}*}');
  })
  .controller('PezAuthController', function($scope, $http, $timeout) {
    var myData = {};
    var pauth = this;
    var restUriBase = "/v1/auth/api-key";
    var meUri = "/me";
    $scope.myApiKey = "You don't have a key yet";

    $timeout(function () {  
      callMeUsingVerb($http.get, meUri);
    }, 1);

    function getRestUri() {
      return [restUriBase, $scope.myEmail].join("/")
    }

    pauth.createorg = function() {
      console.log("not implemented yet");
    };
   
    pauth.create = function() {
      callAPIUsingVerb($http.put, getRestUri())
    };
 
    pauth.remove = function() {
      callAPIUsingVerb($http.delete, getRestUri())
    };   

    function callMeUsingVerb(verbCaller, uri) {
      var responsePromise = verbCaller(uri);
      responsePromise.success(function(data, status, headers, config) {
          $scope.myName = data.Payload.displayName;
          $scope.myEmail = data.Payload.emails[0].value
          callAPIUsingVerb($http.get, getRestUri());
      });
      
      responsePromise.error(function(data, status, headers, config) {
          console.log("AJAX failed!");
      });
    }

    function callAPIUsingVerb(verbCaller, uri) {
      var responsePromise = verbCaller(uri);
      responsePromise.success(function(data, status, headers, config) {
          $scope.myData = data;
          $scope.myApiKey = data.APIKey;
      });
      
      responsePromise.error(function(data, status, headers, config) {
          console.log("AJAX failed!");
      });
    }
  });
