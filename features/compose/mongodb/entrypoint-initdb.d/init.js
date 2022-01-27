var databases = [
	{
		name: "census",
		collections: ["filters", "filterOutputs"]
	}
];

for (database of databases) {
	db = db.getSiblingDB(database.name);

	for (collection of database.collections){
		db.collection(collection).drop();
		db.createCollection(collection);
	}
}
