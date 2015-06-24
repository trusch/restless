var Model = {
  __uid: '',
  __rev: 0,
  __data: {},
  __schema: '',

  initFromData: function(data){
    this.__data = data;
    return this;
  },

  initFromUID: function(uid){
    if(typeof this._preInit === 'function'){
      this._preInit();
    }
    if(typeof this._doInit === 'function'){
      this._doInit(uid); //impl by subclass
    }
    if(typeof this._postInit === 'function'){
      this._postInit();
    }
    return this;
  },

  remove: function(){
    var res = false;
    if(typeof this._preRemove === 'function'){
      this._preRemove();
    }
    if(typeof this._doRemove === 'function'){
      res = this._doRemove(); //impl by subclass
    }
    if(typeof this._postRemove === 'function'){
      this._postRemove();
    }
    return res;
  },

  commit: function(){
    var res = false;
    if(typeof this._preCommit === 'function'){
      this._preCommit();
    }
    if(typeof this._doCommit === 'function'){
      res = this._doCommit(); //impl by subclass
    }
    if(typeof this._postCommit === 'function'){
      this._postCommit();
    }
    return res;
  }
};
