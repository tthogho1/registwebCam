import boto3
from boto3.session import Session

session = Session(profile_name='myregion')

s3 = session.resource('s3')
bucket = s3.Bucket('bucket4image')

for obj in bucket.objects.all():
    bucket.copy(
        {
            'Bucket': 'bucket4image',
            'Key': obj.key
        },
        obj.key,
        ExtraArgs={
            'ContentType': 'image/jpeg',
            'MetadataDirective': 'REPLACE'
        }
    )
    print(f"Updated metadata for {obj.key}") 