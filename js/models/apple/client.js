/*global ClientModel*/
/*global CommonApple*/
/*global InstanceBuilder*/

var Apple = {
  _preInit: function(){
    console.log('preInit');
  },
  _preRemove: function(){
    console.log('preRemove');
  },
  _preCommit: function(){
    console.log('preCommit');
  },
  _postInit: function(){
    console.log('postInit');
  },
  _postRemove: function(){
    console.log('postRemove');
  },
  _postCommit: function(){
    console.log('postCommit');
  }
};

Apple = InstanceBuilder.extend(ClientModel, CommonApple, Apple);
