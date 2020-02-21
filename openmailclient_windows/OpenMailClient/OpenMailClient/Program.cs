using System;
using System.Diagnostics;
using System.Net.Mail;


namespace OpenMailClient
{
    class Program
    {
        static void Main(string[] args)
        {
            var userFolder = System.IO.Path.Combine(
                Environment.GetFolderPath(Environment.SpecialFolder.UserProfile),
                "m4backup");
            System.IO.Directory.CreateDirectory(userFolder);
            var filesToDelete = System.IO.Directory.GetFiles(userFolder);
            foreach (string emlFile in filesToDelete)
            {
                System.IO.File.Delete(emlFile);
            }

            var mail = new MailMessage(new MailAddress("change@me.net"),new MailAddress("change@me.net"));
            foreach (string attach in args)
            {
                mail.Attachments.Add(new Attachment(attach));
            }
            var client = new SmtpClient();
            client.UseDefaultCredentials = true;
            client.DeliveryMethod = SmtpDeliveryMethod.SpecifiedPickupDirectory;
            client.PickupDirectoryLocation = userFolder;
            try
            {
                client.Send(mail);
            }catch(Exception e)
            {
                Console.WriteLine(e.ToString());
            }
            var files = System.IO.Directory.GetFiles(userFolder);
            foreach (string emlFile in files)
            {
                var process = new Process();
                process.StartInfo.FileName = "explorer.exe";
                process.StartInfo.Arguments = emlFile;
                process.Start();
            }
        }
    }
}
