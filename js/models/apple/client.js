/**
 * Apple Model: A simple model example
 */
restless.Model('Apple', CommonApple, class {
  onPreInit(){
    console.debug('preInit');
  }
  onPreRemove(){
    console.debug('preRemove');
  }
  onPreCommit(){
    console.debug('preCommit');
  }
  onPostInit(){
    console.debug('postInit');
  }
  onPostRemove(){
    console.debug('postRemove');
  }
  onPostCommit(){
    console.debug('postCommit');
  }
});

