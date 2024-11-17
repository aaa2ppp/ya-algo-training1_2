using System;
using System.IO;
using System.Linq;
using System.Collections.Generic;

class Program
{
    static void Main()
    {
        // Parse the first line with number of chairs (n) and height (H) of Vasya
        var firstLine = Console.ReadLine()
                    .Trim()
                    .Split(new[] { ' ' }, StringSplitOptions.RemoveEmptyEntries)
                    .Select(int.Parse)
                    .ToArray();
        int n = firstLine[0]; // Number of chairs
        int H = firstLine[1]; // Vasya's height

        // Parse the second line for chair heights, handling extra spaces
        var heights = Console.ReadLine()
                    .Trim().Split(new[] { ' ' }, StringSplitOptions.RemoveEmptyEntries)
                    .Select(int.Parse)
                    .ToArray();

        // Parse the third line for chair widths, handling extra spaces
        var widths = Console.ReadLine()
                    .Trim().Split(new[] { ' ' }, StringSplitOptions.RemoveEmptyEntries)
                    .Select(int.Parse)
                    .ToArray();
 
        
        // Step 2: Create list of chairs as tuples (height, width) and sort by height, then width
        var chairs = new (long height, long width)[n];
        for (int i = 0; i < n; i++)
        {
            //Early exit if have a chair with w>=H
            if (widths[i] >= H)
            {
                Console.WriteLine(0);
                return;
            }
            chairs[i] = (heights[i], widths[i]);
        }
        chairs = chairs.OrderBy(chair => chair.height).ToArray();
        //Console.WriteLine(string.Join(" ", chairs));
        // Step 3: Sliding window setup
        int left = 0;
        long currentWidthSum = chairs[0].width;
        long minDiscomfort = long.MaxValue;

        // One deque to track maximum consecutive height differences in the current window
        var discomfortDeque = new LinkedList<long>();  // Stores maximum consecutive height differences

        // Step 4: Slide the window across the chairs array
        for (int right = 1; right < n; right++)
        {
            // Expand the window by adding the right chair's width
            currentWidthSum += chairs[right].width;

            // Calculate consecutive height difference and maintain discomfortDeque
                long heightDifference = chairs[right].height - chairs[right - 1].height;

                // Keep deque in decreasing order to maintain the maximum at the front
                while (discomfortDeque.Count > 0 && discomfortDeque.Last() < heightDifference)
                {
                    discomfortDeque.RemoveLast();
                }
                discomfortDeque.AddLast(heightDifference);


            // Shrink the window from the left while width >= H
            while (currentWidthSum >= H)
            {
                // Check the maximum consecutive height difference in the current valid window
                if (discomfortDeque.Count > 0)
                {
                    // Update minimum discomfort if the current max difference is smaller
                    minDiscomfort = Math.Min(minDiscomfort, discomfortDeque.First());
                    if (chairs[left + 1].height - chairs[left].height == discomfortDeque.First())
                    {
                        discomfortDeque.RemoveFirst();
                    }
                }

                // Slide the left side of the window
                currentWidthSum -= chairs[left].width;
                left++;
            }
        }

        // Output the minimum discomfort found
        Console.WriteLine(minDiscomfort);
    }
}
