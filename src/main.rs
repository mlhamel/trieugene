use faktory::ConsumerBuilder;
use std::io;

fn main() {
    let mut c = ConsumerBuilder::default();
    c.register("foobar", |job| -> io::Result<()> {
        println!("{:?}", job);
        Ok(())
    });

    let mut c = c.connect(None).unwrap();
    if let Err(e) = c.run(&["default"]) {
        println!("worker failed: {}", e);
    }
}
