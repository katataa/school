import 'package:flutter/material.dart';
import 'package:hired/models/job.dart';
import 'package:provider/provider.dart';
import '../providers/job_provider.dart';
import '../screens/create_job_screen.dart';
import '../widgets/job_card.dart';
import '../widgets/theme_toggle_button.dart';

class JobListScreen extends StatefulWidget {
  final VoidCallback toggleTheme;
  final bool isDarkMode;

  const JobListScreen({
    super.key,
    required this.toggleTheme,
    required this.isDarkMode,
  });

  @override
  JobListScreenState createState() => JobListScreenState();
}

class JobListScreenState extends State<JobListScreen> {
  String selectedFilter = 'All Jobs';

  @override
  void initState() {
    super.initState();
    Provider.of<JobProvider>(context, listen: false).fetchJobs();
  }

  @override
  Widget build(BuildContext context) {
    final jobProvider = Provider.of<JobProvider>(context);
    final jobs = _filterJobs(jobProvider.jobs);

    return Scaffold(
      appBar: AppBar(
        title: const Text('Job List'),
        actions: [
          Padding(
            padding: const EdgeInsets.only(right: 16.0),
            child: DropdownButtonHideUnderline(
              child: DropdownButton<String>(
                value: selectedFilter,
                style: TextStyle(
                  color:
                      widget.isDarkMode
                          ? Colors.white
                          : Colors.black,
                ),
                borderRadius: BorderRadius.circular(12),
                onChanged:
                    (newValue) => setState(() => selectedFilter = newValue!),
                items:
                    [
                          'All Jobs',
                          'Todo',
                          'In Progress',
                          'Complete',
                          'Debtors',
                          'Cancelled',
                        ]
                        .map(
                          (filter) => DropdownMenuItem(
                            value: filter,
                            child: Text(filter),
                          ),
                        )
                        .toList(),
              ),
            ),
          ),
          IconButton(
  icon: const Icon(Icons.delete_forever),
  tooltip: 'Delete All Jobs',
  onPressed: () {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete All Jobs'),
        content: const Text('Are you sure you want to delete all jobs?'),
        actions: [
          TextButton(onPressed: () => Navigator.pop(context), child: const Text('Cancel')),
          TextButton(
            onPressed: () {
              Provider.of<JobProvider>(context, listen: false).deleteAllJobs();
              Navigator.pop(context);
            },
            child: const Text('Delete', style: TextStyle(color: Colors.red)),
          ),
        ],
      ),
    );
  },
)

        ],
        
      ),
      
      body: AnimatedSwitcher(
        duration: const Duration(
          milliseconds: 500,
        ),
        transitionBuilder: (Widget child, Animation<double> animation) {
          return FadeTransition(
            opacity: animation,
            child: child,
          );
        },
        child: Stack(
          key: ValueKey<bool>(
            widget.isDarkMode,
          ),
          children: [
            jobs.isEmpty
                ? const Center(
                  child: Text(
                    'No jobs in this field',
                    style: TextStyle(fontSize: 18, color: Colors.grey),
                  ),
                )
                : ListView.builder(
                  itemCount: jobs.length,
                  itemBuilder: (context, index) {
                    final job = jobs[index];
                    return AnimatedOpacity(
                      opacity: 1,
                      duration: Duration(milliseconds: 300 + (index * 100)),
                      child: JobCard(
                        job: job,
                        isDarkMode:
                            widget.isDarkMode,
                      ),
                    );
                  },
                ),
            AnimatedIconButton(
              isDarkMode: widget.isDarkMode,
              onToggle: widget.toggleTheme,
            ),
          ],
        ),
      ),
      floatingActionButton: FloatingActionButton(
  heroTag: 'addJob',
  onPressed: () {
    Navigator.push(
      context,
      MaterialPageRoute(builder: (context) => const CreateJobScreen()),
    );
  },
  tooltip: 'Add Job',
  child: const Icon(Icons.add),
),

      
    );
  }

  List<Job> _filterJobs(List<Job> jobs) {
    print('Selected Filter: $selectedFilter');

    if (selectedFilter == 'All Jobs') {
      return jobs;
    }

    if (selectedFilter == 'Debtors') {
      return jobs.where((job) => job.cost > job.amountPaid).toList();
    }

    for (var job in jobs) {
      print('Job Status: "${job.status.trim().toLowerCase()}"');
      print('Comparing: "${selectedFilter.trim().toLowerCase()}"');
    }

    return jobs
        .where(
          (job) =>
              job.status.trim().toLowerCase() ==
              selectedFilter.trim().toLowerCase(),
        )
        .toList();
  }
}
