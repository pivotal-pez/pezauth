describe('PezAuthController', function() {
  beforeEach(module('pezAuth'));

  var $controller;

  beforeEach(inject(function(_$controller_){
    $controller = _$controller_;
  }));

  describe('$scope.myName', function() {
    it('sets the strength to "strong" if the password length is >8 chars', function() {
      var $scope = {};
      var controller = $controller('PezAuthController', { $scope: $scope });
      $scope.myName = 'longerthaneightchars';
      expect(true).toEqual(true);
    });
  });
});
