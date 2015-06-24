var InstanceBuilder = {
  extend: function(){
    if(arguments.length < 2){
      throw 'specify at least 2 arguments';
    }
    var result = {};
    for(var i=0;i<arguments.length;i++){
      for(var key in arguments[i]){
        if(arguments[i].hasOwnProperty(key)){
          result[key] = arguments[i][key];
        }
      }
    }
    return result;
  }
};
