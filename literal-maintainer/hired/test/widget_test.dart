import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:hired/main.dart';
import 'package:hired/screens/job_list_screen.dart';

void main() {
  testWidgets('Job List Screen loads correctly', (WidgetTester tester) async {
    // Build the app and trigger a frame.
    await tester.pumpWidget(const MyApp());

    // Verify that the Job List screen is shown.
    expect(find.byType(JobListScreen), findsOneWidget);
    expect(find.text('Job List'), findsOneWidget);
  });

  testWidgets('Displays "No jobs in this field" when no jobs exist', (
    WidgetTester tester,
  ) async {
    await tester.pumpWidget(const MyApp());

    // Check if the default empty state message appears when no jobs exist.
    expect(find.text('No jobs in this field'), findsOneWidget);
  });

  testWidgets('Job status filter dropdown works', (WidgetTester tester) async {
    await tester.pumpWidget(const MyApp());

    // Find the dropdown button
    final dropdownFinder = find.byType(DropdownButton<String>);
    expect(dropdownFinder, findsOneWidget);

    // Open the dropdown
    await tester.tap(dropdownFinder);
    await tester.pump();

    // Select "Completed"
    await tester.tap(find.text('Completed').last);
    await tester.pump();

    // Verify that the selected filter text updates
    expect(find.text('Completed'), findsOneWidget);
  });

  bool jobExists(String jobTitle) {
    return find.text(jobTitle).evaluate().isNotEmpty;
  }

  testWidgets('User can create a new job', (WidgetTester tester) async {
    await tester.pumpWidget(const MyApp());

    // Tap the floating action button to open the Create Job screen
    await tester.tap(find.byIcon(Icons.add));
    await tester.pumpAndSettle();

    // Fill out the form fields
    await tester.enterText(find.byLabelText('Title *'), 'Fix sink');
    await tester.enterText(find.byLabelText('Contact *'), 'John Doe');
    await tester.enterText(find.byLabelText('Address *'), '123 Main St');
    await tester.enterText(find.byLabelText('Cost *'), '50.00');
    await tester.enterText(find.byLabelText('Amount Paid *'), '20.00');

    // Tap the Save button
    await tester.tap(find.text('Save Job'));
    await tester.pumpAndSettle();

    // Verify the new job appears in the job list
    expect(find.text('Fix sink'), findsOneWidget);
  });

  testWidgets('User can edit an existing job', (WidgetTester tester) async {
    await tester.pumpWidget(const MyApp());

    // Ensure at least one job exists
    if (!jobExists('Fix sink')) {
      await tester.tap(find.byIcon(Icons.add));
      await tester.pumpAndSettle();
      await tester.enterText(find.byLabelText('Title *'), 'Fix sink');
      await tester.enterText(find.byLabelText('Contact *'), 'John Doe');
      await tester.enterText(find.byLabelText('Address *'), '123 Main St');
      await tester.enterText(find.byLabelText('Cost *'), '50.00');
      await tester.enterText(find.byLabelText('Amount Paid *'), '20.00');
      await tester.tap(find.text('Save Job'));
      await tester.pumpAndSettle();
    }

    // Tap on the job to open details
    await tester.tap(find.text('Fix sink'));
    await tester.pumpAndSettle();

    // Edit the title
    await tester.enterText(find.byLabelText('Title *'), 'Fix sink - Updated');
    await tester.tap(find.text('Save Job'));
    await tester.pumpAndSettle();

    // Verify the job title is updated
    expect(find.text('Fix sink - Updated'), findsOneWidget);
  });

  testWidgets('User can delete a job', (WidgetTester tester) async {
    await tester.pumpWidget(const MyApp());

    // Ensure at least one job exists
    if (!jobExists('Fix sink')) {
      await tester.tap(find.byIcon(Icons.add));
      await tester.pumpAndSettle();
      await tester.enterText(find.byLabelText('Title *'), 'Fix sink');
      await tester.enterText(find.byLabelText('Contact *'), 'John Doe');
      await tester.enterText(find.byLabelText('Address *'), '123 Main St');
      await tester.enterText(find.byLabelText('Cost *'), '50.00');
      await tester.enterText(find.byLabelText('Amount Paid *'), '20.00');
      await tester.tap(find.text('Save Job'));
      await tester.pumpAndSettle();
    }

    // Find the job and delete it
    await tester.tap(find.byIcon(Icons.delete).first);
    await tester.pumpAndSettle();

    // Verify the job is deleted
    expect(find.text('Fix sink'), findsNothing);
  });
}

extension on CommonFinders {
  Finder byLabelText(String s) {
    return find.widgetWithText(TextField, s);
  }
}
