import 'dart:io' show Platform;
import 'package:sqflite/sqflite.dart';
import 'package:path/path.dart';
import '../models/job.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class DatabaseService {
  static Database? _database;
  static final DatabaseService instance = DatabaseService._init();

  Database get database {
    if (_database == null) {
      throw Exception("Database not initialized");
    }
    return _database!;
  }

  DatabaseService._init();

  final _secureStorage = const FlutterSecureStorage();

  Future<void> initializeDatabase() async {
    final databasePath = await getDatabasesPath();
    final dbFullPath = join(databasePath, 'jobs.db');
    print('Database path: $databasePath');

    // TEMPORARY RESET – delete old DB to avoid NOT NULL crash
     // await deleteDatabase(dbFullPath);
 //    print('🗑️ Deleted old database');

    _database = await openDatabase(
      dbFullPath,
      version: 3,
      onCreate: (db, version) {
        return db.execute('''
          CREATE TABLE jobs(
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            status TEXT NOT NULL,
            info TEXT,
            contact TEXT,
            address TEXT,
            description TEXT,
            cost REAL NOT NULL,
            amountPaid REAL NOT NULL,
            images TEXT
          )
        ''');
      },
    );
  }

  Future<void> saveSecureData(String key, String value) async {
    if (Platform.isMacOS) return; // Skip on macOS
    await _secureStorage.write(key: key, value: value);
  }

  Future<String?> getSecureData(String key) async {
    if (Platform.isMacOS) return null;
    return await _secureStorage.read(key: key);
  }

  Future<List<Job>> getJobs() async {
  final List<Map<String, dynamic>> maps = await _database!.query('jobs');

  return Future.wait(maps.map((map) async {
    final id = map['id'];
    final secureContact = await getSecureData('contact_$id');
    final secureAddress = await getSecureData('address_$id');

    final job = Job.fromMap(map,
      contact: secureContact,
      address: secureAddress,
    );

    // uncomment to see confirmation on data being encryted
    //print('🔐 Loaded secure: contact: ${job.contact}, address: ${job.address}');
    return job;
  }));
}

  Future<void> insertJob(Job job) async {
    print('Inserting job: ${job.title}');

    // Prepare map without contact & address
    final jobMap = job.toMap();
    final contact = job.contact;
    final address = job.address;

    jobMap.remove('contact');
    jobMap.remove('address');

    final id = await _database!.insert(
      'jobs',
      jobMap,
      conflictAlgorithm: ConflictAlgorithm.replace,
    );
    job.id = id;

    await saveSecureData('contact_$id', contact);
    await saveSecureData('address_$id', address);

    // uncomment to see confirmation on data being encryted
    //print('✅ Saved securely: contact_$id = $contact | address_$id = $address');
  }

  Future<void> updateJob(Job job) async {
    await _database!.update(
      'jobs',
      job.toMap(),
      where: 'id = ?',
      whereArgs: [job.id],
    );
  }

  Future<void> deleteJob(int id) async {
    await _database!.delete('jobs', where: 'id = ?', whereArgs: [id]);
  }
}
