import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:provider/provider.dart';
import '../models/job.dart';
import '../providers/job_provider.dart';

class JobDetailsScreen extends StatefulWidget {
  final Job job;
  const JobDetailsScreen({super.key, required this.job});

  @override
  JobDetailsScreenState createState() => JobDetailsScreenState();
}

class JobDetailsScreenState extends State<JobDetailsScreen> {
  final _formKey = GlobalKey<FormState>();
  late String title, status, info, contact, address, description;
  late double cost, amountPaid;

  @override
  void initState() {
    super.initState();
    title = widget.job.title;
    status = widget.job.status;
    info = widget.job.info;
    contact = widget.job.contact;
    address = widget.job.address;
    description = widget.job.description;
    cost = widget.job.cost;
    amountPaid = widget.job.amountPaid;
  }

  void _saveChanges() {
    if (_formKey.currentState!.validate()) {
      final updatedJob = Job(
        id: widget.job.id,
        title: title,
        status: status,
        info: info,
        contact: contact,
        address: address,
        description: description,
        cost: cost,
        amountPaid: amountPaid,
        images: widget.job.images,
      );

      Provider.of<JobProvider>(context, listen: false).updateJob(updatedJob);
      Navigator.pop(context);
    }
  }

  void _deleteJob() {
    Provider.of<JobProvider>(context, listen: false).deleteJob(widget.job.id!);
    Navigator.pop(context);
  }

  void _copyToClipboard() {
    final jobText = '''
Title: $title
Status: $status
Contact: $contact
Address: $address
Description: $description
Cost: ${cost.toStringAsFixed(2)}
Amount Paid: ${amountPaid.toStringAsFixed(2)}
Info: ${info.isNotEmpty ? info : "No additional info"}
''';

    Clipboard.setData(ClipboardData(text: jobText));
    ScaffoldMessenger.of(context).showSnackBar(
      const SnackBar(content: Text('Job details copied to clipboard')),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Job Details')),
      body: Padding(
        padding: const EdgeInsets.all(20),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              if (widget.job.images.isNotEmpty)
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    const Text("Job Images", style: TextStyle(fontWeight: FontWeight.bold)),
                    const SizedBox(height: 10),
                    Wrap(
                      spacing: 10,
                      runSpacing: 10,
                      children: widget.job.images
    .where((imgPath) => imgPath.isNotEmpty)
    .map((imgPath) => ClipRRect(
          borderRadius: BorderRadius.circular(12),
          child: Image.file(
            File(imgPath),
            width: 100,
            height: 100,
            fit: BoxFit.cover,
          ),
        ))
    .toList()

                    ),
                    const SizedBox(height: 20),
                  ],
                ),
              _buildField('Title', Icons.work, title, (val) => title = val),
              _buildStatusDropdown(),
              _buildField('Contact', Icons.person, contact, (val) => contact = val),
              _buildField('Address', Icons.location_on, address, (val) => address = val),
              _buildField('Cost', Icons.money, cost.toString(), (val) => cost = double.tryParse(val) ?? 0),
              _buildField('Amount Paid', Icons.attach_money, amountPaid.toString(), (val) => amountPaid = double.tryParse(val) ?? 0),
              _buildField('Description', Icons.description, description, (val) => description = val),
              const SizedBox(height: 20),
              ElevatedButton(
                onPressed: _saveChanges,
                child: const Text('Save Changes'),
              ),
              const SizedBox(height: 10),
              ElevatedButton(
                onPressed: _deleteJob,
                style: ElevatedButton.styleFrom(backgroundColor: Colors.red),
                child: const Text('Delete Job'),
              ),
              const SizedBox(height: 10),
              ElevatedButton.icon(
                onPressed: _copyToClipboard,
                icon: const Icon(Icons.copy),
                label: const Text('Copy to Clipboard'),
                style: ElevatedButton.styleFrom(
                  backgroundColor: Colors.indigo,
                  foregroundColor: Colors.white,
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildField(String label, IconData icon, String initial, Function(String) onChanged) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: TextFormField(
        initialValue: initial,
        decoration: InputDecoration(labelText: label, prefixIcon: Icon(icon)),
        onChanged: onChanged,
      ),
    );
  }

  Widget _buildStatusDropdown() {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: DropdownButtonFormField<String>(
        value: status,
        decoration: const InputDecoration(
          labelText: 'Status',
          prefixIcon: Icon(Icons.check_circle),
        ),
        onChanged: (newValue) => setState(() => status = newValue!),
        items: [
          'Todo',
          'In Progress',
          'Complete',
          'Cancelled',
        ].map((status) => DropdownMenuItem(
              value: status,
              child: Text(status),
            )).toList(),
      ),
    );
  }
}
