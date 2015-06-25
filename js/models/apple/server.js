/*global ServerModel*/
/*global CommonApple*/

class AppleDef {
  _preInit(){
    console.log('preInit');
  }
  _preRemove(){
    console.log('preRemove');
  }
  _preCommit(){
    console.log('preCommit');
  }
  _postInit(){
    console.log('postInit');
  }
  _postRemove(){
    console.log('postRemove');
  }
  _postCommit(){
    console.log('postCommit');
  }
};

class Apple extends aggregation(ServerModel, CommonApple, AppleDef) { }