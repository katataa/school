import 'dart:io';
import 'package:flutter/material.dart';
import 'package:image_picker/image_picker.dart';
import 'package:provider/provider.dart';
import '../models/job.dart';
import '../providers/job_provider.dart';

String info = '';

class CreateJobScreen extends StatefulWidget {
  const CreateJobScreen({super.key});

  @override
  CreateJobScreenState createState() => CreateJobScreenState();
}

class CreateJobScreenState extends State<CreateJobScreen> {
  final _formKey = GlobalKey<FormState>();
  final _titleController = TextEditingController();

  String status = 'Todo';
  String contact = '';
  String address = '';
  double cost = 0.0;
  double amountPaid = 0.0;
  List<String> images = [];

  @override
  void dispose() {
    _titleController.dispose();
    super.dispose();
  }

  void _saveJob() {
    if (_formKey.currentState!.validate()) {
      final newJob = Job(
        title: _titleController.text,
        status: status,
        contact: contact,
        address: address,
        description: info,
        cost: cost,
        amountPaid: amountPaid,
        images: images,
        info: info,
      );
      Provider.of<JobProvider>(context, listen: false).addJob(newJob);
      Navigator.pop(context);
    }
  }

  Future<void> _pickImage(ImageSource source) async {
    final pickedFile = await ImagePicker().pickImage(source: source);
    if (pickedFile != null && images.length < 5) {
      setState(() => images.add(pickedFile.path));
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: const Text('Create Job')),
      body: Padding(
        padding: const EdgeInsets.all(20),
        child: Form(
          key: _formKey,
          child: ListView(
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 8),
                child: TextFormField(
                  controller: _titleController,
                  decoration: const InputDecoration(
                    labelText: 'Title *',
                    prefixIcon: Icon(Icons.work),
                  ),
                  validator: (value) => value!.isEmpty ? 'Title is required' : null,
                ),
              ),
              _buildDropdown(),
              _buildTextField('Contact *', Icons.person, (val) => contact = val),
              _buildTextField('Address *', Icons.location_on, (val) => address = val),
              _buildTextField(
                'Cost *',
                Icons.monetization_on,
                (val) => cost = double.tryParse(val) ?? 0.0,
                keyboardType: TextInputType.number,
              ),
              _buildTextField(
                'Amount Paid *',
                Icons.attach_money,
                (val) => amountPaid = double.tryParse(val) ?? 0.0,
                keyboardType: TextInputType.number,
              ),
              const SizedBox(height: 10),
              TextFormField(
                maxLines: 5,
                decoration: const InputDecoration(
                  labelText: 'Description (Optional)',
                  alignLabelWithHint: true,
                  prefixIcon: Icon(Icons.description),
                ),
                onChanged: (val) => setState(() => info = val),
              ),
              const SizedBox(height: 20),
              const Text(
                'Upload Images (Max 5)',
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              Wrap(
                spacing: 10,
                runSpacing: 10,
                children: [
                  ...images.map(
                    (imgPath) => ClipRRect(
                      borderRadius: BorderRadius.circular(12),
                      child: Image.file(
                        File(imgPath),
                        width: 100,
                        height: 100,
                        fit: BoxFit.cover,
                      ),
                    ),
                  ),
                  if (images.length < 5)
                    InkWell(
                      onTap: () => _showImagePicker(context),
                      child: Container(
                        width: 100,
                        height: 100,
                        decoration: BoxDecoration(
                          borderRadius: BorderRadius.circular(12),
                          color: Colors.grey[200],
                        ),
                        child: const Icon(Icons.add_a_photo, size: 30),
                      ),
                    ),
                ],
              ),
              const SizedBox(height: 30),
              ElevatedButton(
                onPressed: _saveJob,
                child: const Text('Save Job'),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildTextField(
    String label,
    IconData icon,
    Function(String) onChanged, {
    TextInputType keyboardType = TextInputType.text,
  }) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: TextFormField(
        keyboardType: keyboardType,
        decoration: InputDecoration(labelText: label, prefixIcon: Icon(icon)),
        validator: (value) => value!.isEmpty ? '$label is required' : null,
        onChanged: onChanged,
      ),
    );
  }

  Widget _buildDropdown() {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: DropdownButtonFormField<String>(
        value: status,
        decoration: const InputDecoration(
          labelText: 'Status *',
          prefixIcon: Icon(Icons.check_circle),
        ),
        onChanged: (newValue) => setState(() => status = newValue!),
        items: ['Todo', 'In Progress', 'Complete', 'Cancelled']
            .map((status) => DropdownMenuItem(value: status, child: Text(status)))
            .toList(),
      ),
    );
  }

  void _showImagePicker(BuildContext context) {
    if (!(Platform.isAndroid || Platform.isIOS)) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Image picking only works on Android/iOS for now')),
      );
      return;
    }

    showModalBottomSheet(
      context: context,
      builder: (_) => Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          ListTile(
            leading: const Icon(Icons.photo_library),
            title: const Text('Pick from Gallery'),
            onTap: () {
              Navigator.pop(context);
              _pickImage(ImageSource.gallery);
            },
          ),
        ],
      ),
    );
  }
}
