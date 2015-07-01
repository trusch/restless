/**
 * Apple Model: A simple model example
 */
app.Model('Apple', CommonApple, class {
  
  onPreGet(){
    if(ROLE !== 'admin') throw 'no right!';
  }
  onPostGet(){
    this.set('example.dynamic.requestTime',(new Date()).toString());
  }

  onPreRemove(){
    if(ROLE !== 'admin') throw 'no right!';
  }
  onPostRemove(){
    console.debug(`deleted instance ${this.__uid} of model ${this.__name}`);
  }

  onPrePut(){
    if(ROLE !== 'admin') throw 'no right!';
    this.set('lastModified',Date.now());
  }
  onPostPut(){
    console.debug(`commited instance ${this.__uid} of model ${this.__name}`);
  }
  
  
});

