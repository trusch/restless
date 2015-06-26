class ClientApp {
    constructor(){
        this.__models = {};
    }

    Model(name, ...mixins){
        console.log(`declare server model ${name}...`);
        this.__models[name] = aggregation(ServerModel, ...mixins);
    }

    CreateModel(name, ...args){
        var model = new this.__models[name](...args);
        model.__name = name;
        return model;
    }

}

var app  = new ClientApp();