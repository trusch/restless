/*global InstanceBuilder*/
/*global Model*/

var ClientModel = {
  _doInit: function(uid){
    console.log('ClientModel: init');
  },
  _doRemove: function(){
    console.log('ClientModel: remove');
  },
  _doCommit: function(){
    console.log('ClientModel: commit');
  }
};

ClientModel = InstanceBuilder.extend(Model,ClientModel);

