describe('PezAuthController', function() {
  beforeEach(module('pezAuth'));

  var $controller;

  beforeEach(inject(function(_$controller_){
    $controller = _$controller_;
  }));

  describe('$scope.myName & myEmail', function() {
    it('should be initialized as undefined', function() {
      var $scope = {};
      var controller = $controller('PezAuthController', { $scope: $scope });
      expect($scope.myName).toEqual(undefined);
      expect($scope.myEmail).toEqual(undefined);
    });
  });
});
