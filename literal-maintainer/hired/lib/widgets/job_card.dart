import 'package:flutter/material.dart';
import '../models/job.dart';
import '../screens/job_details_screen.dart';

class JobCard extends StatelessWidget {
  final Job job;
  final bool isDarkMode;

  const JobCard({
    super.key,
    required this.job,
    required this.isDarkMode, 
  });

  Color _statusColor(String status) {
    switch (status) {
      case 'Todo':
        return Colors.grey;
      case 'In Progress':
        return Colors.orange;
      case 'Complete':
        return Colors.green;
      case 'Cancelled':
        return Colors.red;
      default:
        return Colors.blue;
    }
  }

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap:
          () => Navigator.push(
            context,
            PageRouteBuilder(
              pageBuilder: (_, __, ___) => JobDetailsScreen(job: job),
              transitionsBuilder:
                  (_, animation, __, child) =>
                      FadeTransition(opacity: animation, child: child),
            ),
          ),
      child: Container(
        margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        padding: const EdgeInsets.all(16),
        decoration: BoxDecoration(
          color:
              isDarkMode
                  ? Colors.grey[850]
                  : Colors.white,
          borderRadius: BorderRadius.circular(16),
          boxShadow: [
            BoxShadow(
              color:
                  isDarkMode
                      ? Colors.white24
                      : Colors.black12,
              blurRadius: 6,
              offset: const Offset(0, 3),
            ),
          ],
        ),
        child: Row(
          children: [
            CircleAvatar(
              backgroundColor: _statusColor(job.status).withOpacity(0.2),
              child: Icon(Icons.work, color: _statusColor(job.status)),
            ),
            const SizedBox(width: 16),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    job.title,
                    style: TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.w600,
                      color:
                          isDarkMode
                              ? Colors.white
                              : Colors.black,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    job.contact,
                    style: TextStyle(
                      color: isDarkMode ? Colors.white70 : Colors.grey[600],
                    ),
                  ),
                ],
              ),
            ),
            Container(
              padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
              decoration: BoxDecoration(
                color: _statusColor(job.status).withOpacity(0.1),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Text(
                job.status,
                style: TextStyle(
                  fontSize: 12,
                  fontWeight: FontWeight.w500,
                  color: _statusColor(job.status),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
