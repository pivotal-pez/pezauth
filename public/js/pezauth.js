var pezAuth = angular.module('pezAuth', [], function($interpolateProvider) {
      $interpolateProvider.startSymbol('{*{');
      $interpolateProvider.endSymbol('}*}');
  })
  .controller('PezAuthController', function($scope, $http, $timeout, $window) {
    var myData = {};
    var pauth = this;
    var restUriBase = "/v1/auth/api-key";
    var restOrgUriBase = "/v1/org/user";
    var meUri = "/me";
    var messaging = {
      "hasOrgBtn": "View Org Now",
      "createOrgBtn": "Create Your Org Now",
      "noApiKey": "You don't have a key yet"
    };
    var urls = {
      "okta": "http://login.run.pez.pivotal.io/saml/login/alias/login.run.pez.pivotal.io?disco=true"
    };
    

    $timeout(function () {  
      callMeUsingVerb($http.get, meUri);
    }, 1);

    function getRestUri() {
      return [restUriBase, $scope.myEmail].join("/")
    }

    function getOrgRestUri() {
      return [restOrgUriBase, $scope.myEmail].join("/")
    }

    pauth.getorg = function() {
      console.log(getOrgRestUri())
      callOrgUsingVerb($http.get, getOrgRestUri())
    };

    pauth.createorg = function() {
      //console.log(getOrgRestUri())
      //callOrgUsingVerb($http.put, getOrgRestUri())
      //console.log("not implemented yet");
      $window.location.href = urls.okta;
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
          pauth.getorg()
      });
    }

    function callOrgUsingVerb(verbCaller, uri) {
      var responsePromise = verbCaller(uri);
      responsePromise.success(function(data, status, headers, config) {
        console.log(data);
        $scope.orgButtonText = messaging.hasOrgBtn;
      });
      
      responsePromise.error(function(data, status, headers, config) {
          $scope.orgButtonText = messaging.createOrgBtn;
          
          if(status === 403 && verbCaller === $http.get) {
            console.log(data.ErrorMsg);
          }
      });
    }

    function callAPIUsingVerb(verbCaller, uri) {
      var responsePromise = verbCaller(uri);
      responsePromise.success(function(data, status, headers, config) {
          $scope.myData = data;
          $scope.myApiKey = data.APIKey;
      });
      
      responsePromise.error(function(data, status, headers, config) {
        $scope.myApiKey = messaging.noApiKey;
      });
    }
  });
