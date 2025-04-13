import 'package:flutter/material.dart';

class AnimatedIconButton extends StatefulWidget {
  final bool isDarkMode;
  final VoidCallback onToggle;

  const AnimatedIconButton({
    super.key,
    required this.isDarkMode,
    required this.onToggle,
  });

  @override
  State<AnimatedIconButton> createState() => _AnimatedIconButtonState();
}

class _AnimatedIconButtonState extends State<AnimatedIconButton>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _rotation;

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(
      duration: const Duration(milliseconds: 500),
      vsync: this,
    );
    _rotation = Tween<double>(
      begin: 0,
      end: 1,
    ).animate(CurvedAnimation(parent: _controller, curve: Curves.easeInOut));
  }

  @override
  void didUpdateWidget(covariant AnimatedIconButton oldWidget) {
    super.didUpdateWidget(oldWidget);
    if (widget.isDarkMode != oldWidget.isDarkMode) {
      _controller.forward(from: 0); // Reset and play the animation
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Positioned(
      bottom: 30, // Adjust the bottom position (above the FAB)
      left: 16, // Position the button on the left
      child: GestureDetector(
        onTap: () {
          widget.onToggle(); // Toggle the theme
          _controller.forward(
            from: 0,
          ); // Trigger the animation when the button is clicked
        },
        child: RotationTransition(
          turns: _rotation, // Apply the rotation animation
          child: CircleAvatar(
            backgroundColor: Theme.of(context).colorScheme.primaryContainer,
            child: Icon(
              widget.isDarkMode ? Icons.nightlight_round : Icons.wb_sunny,
              color: Theme.of(context).colorScheme.onPrimaryContainer,
            ),
          ),
        ),
      ),
    );
  }
}
