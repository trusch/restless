class ServerApp {
    constructor(){
        this.__models = {};
        this.__debugLogs = true;
        this._initDebugLogs();
    }

    _initDebugLogs(){
        let orig = console.log;
        if(this.__debugLogs){
            console.debug = (...args) => console.log('ServerJS>',...args);
        }else{
            console.debug = (...args) => {};
        }
    }

    Model(name, ...mixins){
        console.debug(`declare server model ${name}...`);
        this.__models[name] = aggregation(ServerModel, ...mixins);
    }

    CreateModel(name, ...args){
        var model = new this.__models[name](...args);
        model.__name = name;
        return model;
    }

}

var app  = new ServerApp();