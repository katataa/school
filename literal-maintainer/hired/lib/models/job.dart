class Job {
  int? id;
  String title;
  String status;
  String info;
  String contact;
  String address;
  String description;
  double cost;
  double amountPaid;
  List<String> images;

  Job({
    this.id,
    required this.title,
    required this.status,
    this.info = '',
    required this.contact,
    required this.address,
    this.description = '',
    required this.cost,
    required this.amountPaid,
    this.images = const [],
  });

  Map<String, dynamic> toMap() {
    return {
      'id': id,
      'title': title,
      'status': status,
      'info': info,
      'contact': contact,
      'address': address,
      'description': description,
      'cost': cost,
      'amountPaid': amountPaid,
      'images': images.join(','),
    };
  }

  factory Job.fromMap(Map<String, dynamic> map, {String? contact, String? address}) {
  return Job(
    id: map['id'],
    title: map['title'],
    status: map['status'],
    info: map['info'] ?? '',
    contact: contact ?? '',
    address: address ?? '',
    description: map['description'] ?? '',
    cost: map['cost'],
    amountPaid: map['amountPaid'],
    images: (map['images'] as String?)?.split(',') ?? [],
  );
}

}
