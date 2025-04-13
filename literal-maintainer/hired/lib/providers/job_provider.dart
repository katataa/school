import 'package:flutter/material.dart';
import '../models/job.dart';
import '../services/database_service.dart';

class JobProvider extends ChangeNotifier {
  List<Job> _jobs = [];
  List<Job> get jobs => _jobs;

  Future<void> fetchJobs() async {
    _jobs = await DatabaseService.instance.getJobs();
    for (var job in _jobs) {
      final contact = await DatabaseService.instance.getSecureData('contact_${job.id}');
      final address = await DatabaseService.instance.getSecureData('address_${job.id}');
     // uncomment for confirmation on database encrypting
     // print('🔐 SECURE STORAGE - contact: $contact, address: $address');
    }
    notifyListeners();
  }

  Future<void> addJob(Job job) async {
    await DatabaseService.instance.insertJob(job);
    await fetchJobs();
  }

  Future<void> updateJob(Job job) async {
    await DatabaseService.instance.updateJob(job);
    await fetchJobs();
  }

  Future<void> deleteJob(int id) async {
    await DatabaseService.instance.deleteJob(id);
    await fetchJobs();
  }

  Future<void> deleteAllJobs() async {
    final db = DatabaseService.instance.database;
    await db.delete('jobs');
    await fetchJobs();
  }
}
