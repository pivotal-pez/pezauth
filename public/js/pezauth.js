var pezAuth = angular.module('pezAuth', ['ngDialog'], function($interpolateProvider) {
      $interpolateProvider.startSymbol('{*{');
      $interpolateProvider.endSymbol('}*}');
  })
  .controller('PezAuthController', function($scope, ngDialog, $http, $timeout, $window) {
    $scope.hideCLIExample = true;
    var myData = {};
    var pauth = this;
    var restUriBase = "/v1/auth/api-key";
    var restOrgUriBase = "/v1/org/user";
    var meUri = "/me";
    var messaging = {
      "hasOrgBtn": "View Org Now",
      "createOrgBtn": "Create Your Org Now",
      "noApiKey": "You don't have a key yet",
      "loading": "Loading... Please Wait",
      "oktaSetup": "Get Okta Tile for HeritageCF",
      "invalidUser": "query failed. unable to find matching user guid."
    };
    var urls = {
      "okta": "http://login.run.pez.pivotal.io/saml/login/alias/login.run.pez.pivotal.io?disco=true",
      "oktaHome": "https://pivotal.okta.com/app/UserHome"
    };

    $timeout(function () {
      callMeUsingVerb($http.get, meUri);
    }, 1);

    function getRestUri() {
      return [restUriBase, $scope.myEmail].join("/");
    }

    function getOrgRestUri() {
      return [restOrgUriBase, $scope.myEmail].join("/");
    }

    pauth.getorg = function() {
      console.log(getOrgRestUri());
      getOrgStatus(getOrgRestUri());
    };

    pauth.createorg = function() {

      if ($scope.orgButtonText === messaging.createOrgBtn) {
        createOrg(getOrgRestUri());

      } else if ($scope.orgButtonText === messaging.hasOrgBtn) {
        $window.location.href = urls.okta;

      }  else if ($scope.orgButtonText === messaging.oktaSetup) {
        $window.location.href = urls.oktaHome;
      }
      $scope.orgButtonText = messaging.loading;
    };

    pauth.create = function() {
      callAPIUsingVerb($http.put, getRestUri());
    };

    pauth.remove = function() {
      callAPIUsingVerb($http.delete, getRestUri());
    };

    pauth.sandbox = function() {
      console.log("show sandbox dialog");
      ngDialog.open({
        template: 'templates/sandbox.html',
        className: 'ngdialog-theme-default',
        //controller: 'pac',
        scope: $scope
      });
    };

    pauth.submitSandbox = function() {
        var formData = jQuery.param(
          {
            from: $scope.myEmail,
            name: $scope.myName
          }
        );
        $http({
            method: 'POST',
            url: "/sandbox",
            data: formData,
            headers: {'Content-Type': 'application/x-www-form-urlencoded'}
        }).success(function (data, status, headers, config) {
          jQuery("#sandboxRequestResult")
            .addClass("success")
            .text("Sandbox request is sent, you will be notified through email once the sandbox is created");
          jQuery("#sandboxSubmit").hide();
        }).error(function (data, status, headers, config) {
          jQuery("#sandboxRequestResult")
            .addClass("error")
            .text("The request could not be fullfiled, please contact ask@pivotal.io");
          jQuery("#sandboxSubmit").hide();
        })
    };

    function callMeUsingVerb(verbCaller, uri) {
      var responsePromise = verbCaller(uri);
      responsePromise.success(function(data, status, headers, config) {
          $scope.myName = data.Payload.displayName;
          $scope.myEmail = data.Payload.emails[0].value;
          $scope.displayName = $scope.myName ? $scope.myName : $scope.myEmail
          callAPIUsingVerb($http.get, getRestUri());
          pauth.getorg();
      });
    }

    function createOrg(uri) {
      var responsePromise = $http.put(uri);
      responsePromise.success(function(data, status, headers, config) {
        console.log(data);
        $scope.orgButtonText = messaging.hasOrgBtn;
        $scope.hideCLIExample = false;
      });

      responsePromise.error(function(data, status, headers, config) {

          if(status === 403) {
            console.log(data.ErrorMsg);

            if (messaging.invalidUser == data.ErrorMsg) {
              forwardToOkta = confirm("You have not set up your account in Okta yet. Please head over to Okta and click on the `PEZ HeritageCF` tile.");
            }

            if ( forwardToOkta === true) {
              $window.location.href = urls.oktaHome;

            } else {
              $scope.orgButtonText = messaging.oktaSetup;
            }
          }
      });
    }

    function getOrgStatus(uri) {
      var responsePromise = $http.get(uri);
      responsePromise.success(function(data, status, headers, config) {
        console.log(data);
        $scope.orgButtonText = messaging.hasOrgBtn;
        $scope.hideCLIExample = false;
      });

      responsePromise.error(function(data, status, headers, config) {

          if(status === 403) {
            $scope.orgButtonText = messaging.createOrgBtn;
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
        pauth.create();
      });
    }
  });
